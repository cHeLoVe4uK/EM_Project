export interface RequestLoginUser {
  email: string;
  password: string;
}

export interface RequestRegisterUser {
  email: string;
  password: string;
  username: string;
}

export interface RequestCreateChat {
  name: string;
}

export interface ResponseCreateChat {
  id: string;
}

export interface Message {
  id: string;
  chat_id: string;
  content: string;
  created_at: string;
  author_id: string;
  author_name: string;
  is_edited: boolean;
}

export interface Send {
  content: string;
}