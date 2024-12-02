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
  author: string;
  content: string;
  created_at: string; // Дата в формате ISO
  is_edited: boolean;
}
