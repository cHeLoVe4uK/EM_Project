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

export interface RequestDeleteChat {
  chat_id: string;
}

export interface RequestConnectToChat {
  chat_id: string;
}

export interface ResponseCreateChat {
  id: string;
}

export interface ErrResponse {
  error: string;
  msg: string;
}
