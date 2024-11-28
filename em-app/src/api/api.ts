import apiClient from "./apiClient";
import {
  RequestConnectToChat,
  RequestCreateChat,
  RequestDeleteChat,
  RequestLoginUser,
  RequestRegisterUser,
  ResponseCreateChat,
} from "./types";

// Авторизация
export const login = async (data: RequestLoginUser) => {
  return await apiClient.post("/users/login", data); // Исправлен путь
};

export const register = async (data: RequestRegisterUser) => {
  return await apiClient.post("/users", data); // Исправлен путь
};

// Работа с чатами
export const createChat = async (data: RequestCreateChat): Promise<ResponseCreateChat> => {
  const response = await apiClient.post("/chats", data); // Исправлен путь
  return response.data;
};

export const deleteChat = async (data: RequestDeleteChat) => {
  return await apiClient.delete("/chats", { data }); // Уточнение пути
};

// Подключение к WebSocket
export const connectToChat = (data: RequestConnectToChat): WebSocket => {
  const wsUrl = `${import.meta.env.VITE_WS_BASE_URL || "ws://localhost:8080/api/v1"}/chats/connect`;
  const socket = new WebSocket(`${wsUrl}?chat_id=${data.chat_id}`);
  return socket;
};
