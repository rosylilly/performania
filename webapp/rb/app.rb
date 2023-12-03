require "sinatra/base"
require "sinatra/json"
require "mysql2"
require "mysql2-cs-bind"

class App < Sinatra::Base
  SESSION_USER_ID = 'user_id'

  enable :logging

  set :sessions, true
  set :session_secret, ENV.fetch('PERFORMANIA_SESSION_SECRET', '8132176fac749f180d1edfd008d587ad3f3d7bb010a1ab51d3b318cafea9b205')

  configure :development do
    require "sinatra/reloader"
    register Sinatra::Reloader
  end

  helpers do
    def connect_db
      Mysql2::Client.new(
        host: ENV.fetch('DB_HOST', '127.0.0.1'),
        port: ENV.fetch('DB_PORT', '3306'),
        username: ENV.fetch('DB_USERNAME', 'performania'),
        password: ENV.fetch('DB_PASSWORD', 'performania'),
        database: ENV.fetch('DB_DATABASE', 'performania'),
        encoding: 'utf8mb4',
        reconnect: true,
        database_timezone: :utc,
        cast_booleans: true,
        symbolize_keys: true,
      )
    end

    def db
      Thread.current[:db] ||= connect_db
    end

    def transaction(&block)
      db.query('BEGIN')
      ok = false
      begin
        retval = block.call(db)
        db.query('COMMIT')
        ok = true
        retval
      ensure
        db.query('ROLLBACK') unless ok
      end
    end

    def authorize
      halt 401, { error: 'Unauthorized' }.to_json unless session[SESSION_USER_ID]
      user = db.xquery('SELECT * FROM users WHERE id = ?', session[SESSION_USER_ID]).first
      halt 401, { error: 'Unauthorized' }.to_json unless user

      user
    end
  end

  post '/api/initialize' do
    transaction do |db|
      db.query('TRUNCATE FROM users')
      db.query('TRUNCATE FROM posts')
      db.query('TRUNCATE FROM favorites')
      db.query('TRUNCATE FROM blocks')
      db.query('TRUNCATE FROM notifications')
    end

    status 201
    json({ language: 'ruby' })
  end

  post '/api/finalize' do
    score = params[:score]
    errors = params[:errors]

    status 200
    json({ score:, errors: })
  end

  get '/blob/:id/icon' do
    id = params[:id]
    user = db.xquery('SELECT * FROM users WHERE id = ?', id).first
    halt 404 unless icon

    content_type 'image/png'
    user[:icon]
  end

  get '/blob/:id/cover' do
    id = params[:id]
    user = db.xquery('SELECT * FROM users WHERE id = ?', id).first
    halt 404 unless user

    content_type 'image/png'
    user[:cover]
  end

  get '/blob/:id/photo' do
    id = params[:id]
    post = db.xquery('SELECT * FROM posts WHERE id = ?', id).first
    halt 404 unless post

    content_type 'image/png'
    post[:photo]
  end

  post '/api/users' do
    login = params[:login] || 'anonymous'
    icon = params[:icon][:tempfile].read
    cover = params[:cover][:tempfile].read

    id = transaction do |db|
      if db.xquery('SELECT id FROM users WHERE login = ?', login).first
        halt 409, json({ error: 'login already exists' })
      end

      now = Time.now
      db.xquery(
        'INSERT INTO users (login, icon, cover, created_at, updated_at) VALUES (?, ?, ?, ?, ?)',
        login, icon, cover, now, now
      )
      db.last_id.to_s
    end

    session[SESSION_USER_ID] = id

    status 201
    json({ id: , login: })
  end

  post '/api/session' do
    login = params[:login] || 'anonymous'
    user = db.xquery('SELECT * FROM users WHERE login = ?', login).first
    halt 404 unless user

    session[SESSION_USER_ID] = user[:id]

    status 201
    json({ id: user[:id] })
  end

  post '/api/posts' do
    user = authorize

    user_id = user[:id]
    body = params[:body]
    photo = params[:photo][:tempfile].read rescue nil
    now = Time.now

    id = transaction do |db|
      db.xquery(
        'INSERT INTO posts (user_id, body, photo, created_at, updated_at) VALUES (?, ?, ?, ?)',
        user_id, body, photo, now, now
      )
      db.last_id.to_s
    end

    status 201
    json({ id: id, created_at: now })
  end

  post '/api/users/:id/block' do
    user = authorize

    blocked_user_id = params[:id]

    transaction do |db|
      blocked_user = db.xquery('SELECT * FROM users WHERE id = ?', blocked_user_id).first
      halt 404 unless blocked_user

      block = db.xquery(
        'SELECT * FROM blocks WHERE blocker_id = ? AND blocked_id = ?',
        user[:id], blocked_user[:id]
      ).first
      halt 409 if block

      now = Time.now

      db.xquery(
        'INSERT INTO blocks (blocker_id, blocked_id, created_at, updated_at) VALUES (?, ?, ?, ?)',
        user[:id], blocked_user[:id], now, now
      )
    end

    status 201
    json({})
  end

  post '/api/posts/:id/favorites' do
    user = authorize

    post_id = params[:id]

    transaction do |db|
      post = db.xquery('SELECT * FROM posts WHERE id = ?', post_id).first
      halt 404 unless post

      favorite = db.xquery(
        'SELECT * FROM favorites WHERE user_id = ? AND post_id = ?',
        user[:id], post[:id]
      ).first
      halt 409 if favorite

      now = Time.now

      db.xquery(
        'INSERT INTO favorites (user_id, post_id, created_at, updated_at) VALUES (?, ?, ?, ?)',
        user[:id], post[:id], now, now
      )

      db.xquery(
        'INSERT INTO notifications (post_id, created_at, updated_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE updated_at = ?, read = ?',
        post[:id], now, now, now, false
      )
    end

    status 201
    json({})
  end

  get '/api/posts' do
    user = db.xquery('SELECT * FROM users WHERE id = ?', session[SESSION_USER_ID]).first

    posts = []
    db.xquery('SELECT * FROM posts ORDER BY created_at DESC').each do |post|
      blocked_ids = user ? db.xquery('SELECT * FROM blocks WHERE blocker_id = ?', user[:id]).map { |block| block[:blocked_id] } : []

      json_value = post.slice(:id, :body, :created_at, :updated_at)
      owner = db.xquery('SELECT * FROM users WHERE id = ?', post[:user_id]).first
      json_value[:owner] = owner

      favorites = db.xquery('SELECT * FROM favorites WHERE post_id = ?', post[:id])
      json_value[:favorites] = favorites.size

      next if blocked_ids.include?(post[:user_id])
      posts << json_value
    end
    posts = posts[0..20]

    header 'X-Poll-Interval', '1000'
    json(posts)
  end

  get '/api/notifications' do
    user = authorize

    notifications = []
    db.xquery('SELECT * FROM notifications WHERE user_id = ? AND read = ? ORDER BY updated_at DESC', user[:id], false).each do |notification|
      post = db.xquery('SELECT * FROM posts WHERE id = ?', notification[:post_id]).first
      favorites = db.xquery('SELECT * FROM favorites WHERE post_id = ?', post[:id])
      favorite_users = favorites.map { |favorite| db.xquery('SELECT * FROM users WHERE id = ?', favorite[:user_id]).first }

      notifications << notification.merge(post: post, favorite_users: favorite_users[0..2])
    end

    json(notifications)
  end
end
