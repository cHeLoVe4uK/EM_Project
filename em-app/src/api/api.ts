import apiClient from "./apiClient";
import {
  RequestConnectToChat,
  RequestCreateChat,
  RequestDeleteChat,
  RequestDeleteMsg,
  RequestLoginUser,
  RequestRegisterUser,
  RequestUpdateMsg,
  ResponseCreateChat,
} from "./types";

// Авторизация
export const login = async (data: RequestLoginUser) => {
  return await apiClient.post("/user/login", data);
};

export const register = async (data: RequestRegisterUser) => {
  return await apiClient.post("/user/register", data);
};

export const logout = async () => {
  return await apiClient.post("/user/logout");
};

// Работа с чатами
export const createChat = async (data: RequestCreateChat): Promise<ResponseCreateChat> => {
  const response = await apiClient.post("/chat", data);
  return response.data;
};

export const deleteChat = async (data: RequestDeleteChat) => {
  return await apiClient.delete("/chat", { data });
};

export const connectToChat = async (data: RequestConnectToChat) => {
  return await apiClient.post("/chat/connect", data);
};

// Работа с сообщениями
export const deleteMessage = async (data: RequestDeleteMsg) => {
  return await apiClient.delete("/message", { data });
};

export const updateMessage = async (data: RequestUpdateMsg) => {
  return await apiClient.patch("/message", data);
};
