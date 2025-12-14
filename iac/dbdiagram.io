Table logs {
  id int [primary key]
  token uuid [not null]
  data text [not null]
  service varchar [not null]
  container_name varchar
  container_id varchar
  node varchar
  node_id varchar
  hash varchar [not null, note: 'Hash computed from log data after stripping timestamps']
  created_at timestamp [default: 'now()']
}

Table subscriptions {
  id uuid [primary key]
  chat_id bigint [not null]
  token uuid [not null, note: 'Unique identifier each project sends logs with. May be assigned only once for each chat.']
  created_at timestamp [default: 'now()']
}

Ref logs_sub: subscriptions.token - logs.token

Table chat_settings {
  chat_id bigint [primary key, not null]
  collapse_period int [default: 0, note: 'Timeout period during which no notfications are sent for identical logs']
  mute_until timestamp [default: 0, note: 'Timeout period during which no notifications are sent']
  silence_until timestamp [default: 0, note: 'Timeout period during which notifications are sent silently']
  updated_at timestamp
}

// Many-to-Many <> 
// One-to-One - 
// One-to-Many < 
// Many-to-One >
Ref sub_chat_settings: subscriptions.chat_id - chat_settings.chat_id

Table user_settings {
  user_id bigint [primary key, not null]
  username varchar
  lang int
  updated_at timestamp
}


Table labels {
  chat_id bigint [not null]
  username varchar [not null, note: 'Username is extracted from Entities of the message where the user was mentioned']
  labels varchar[]
  updated_at timestamp
  
  indexes {
    (chat_id, username) [pk] // Composite primary key
  }
}

Ref: labels.chat_id - subscriptions.chat_id

Table commands {
  name varchar [not null]
  user_id bigint [not null]
  chat_id bigint [not null]
  stage int [not null, default: -1]
  data json [default: null]
  updated_at timestamp [not null]

  indexes {
    (user_id, chat_id) [pk]
  }
}

Table last_sent {
  chat_id bigint [not null]
  token uuid [not null]
  hash varchar [not null]
  updated_at tiemstamp [not null]

  indexes {
    (chat_id, token, hash) [pk]
  }
}

Ref: last_sent.chat_id - subscriptions.chat_id

Ref: last_sent.token - subscriptions.token