// Запросы
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
  
  export interface RequestDeleteMsg {
    msg_id: string;
  }
  
  export interface RequestUpdateMsg {
    msg_id: string;
    text: string;
  }
  
  // Ответы
  export interface ResponseCreateChat {
    chat_id: string;
  }
  
  // Ошибки
  export interface ErrResponse {
    error: string;
    msg: string;
  }
  