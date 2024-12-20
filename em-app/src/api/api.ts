import apiClient from "./apiClient";
import {
  RequestCreateChat,
  RequestLoginUser,
  RequestRegisterUser,
  ResponseCreateChat,
} from "./types";
import Cookies from "js-cookie";

export const login = async (data: RequestLoginUser) => {
  const response = await apiClient.post("/users/login", data);
  return response.data;
};

export const register = async (data: RequestRegisterUser) => {
  const response = await apiClient.post("/users", data);
  return response.data;
};

export const createChat = async (data: RequestCreateChat): Promise<ResponseCreateChat> => {
  const response = await apiClient.post("/chats", data);
  return response.data;
};

export const fetchChats = async () => {
  const response = await apiClient.get("/chats");
  return response.data;
};

export const fetchMessages = async (chatId: string) => {
  const response = await apiClient.get(`/chats/${chatId}/messages`);
  return response.data;
};

export const connectToChat = (chatId: string): WebSocket => {
  const token = Cookies.get("access_token");
  const wsUrl = `${import.meta.env.VITE_WS_BASE_URL || "ws://localhost:8080/api/v1"}/chats/${chatId}/connect?token=${token}`;
  return new WebSocket(wsUrl);
};