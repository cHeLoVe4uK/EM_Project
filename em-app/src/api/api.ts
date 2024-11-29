import apiClient from "./apiClient";
import {
  RequestConnectToChat,
  RequestCreateChat,
  RequestDeleteChat,
  RequestLoginUser,
  RequestRegisterUser,
  ResponseCreateChat,
} from "./types";

export const login = async (data: RequestLoginUser) => {
  return await apiClient.post("/users/login", data);
};

export const register = async (data: RequestRegisterUser) => {
  return await apiClient.post("/users", data);
};

export const createChat = async (data: RequestCreateChat): Promise<ResponseCreateChat> => {
  const response = await apiClient.post("/chats", data);
  return response.data;
};

export const deleteChat = async (data: RequestDeleteChat) => {
  return await apiClient.delete(`/chats/${data.chat_id}`);
};

export const connectToChat = (data: RequestConnectToChat): WebSocket => {
  const wsUrl = `${import.meta.env.VITE_WS_BASE_URL || "ws://localhost:8080/api/v1"}/chats/${data.chat_id}/connect`;
  const socket = new WebSocket(wsUrl);
  return socket;
};
