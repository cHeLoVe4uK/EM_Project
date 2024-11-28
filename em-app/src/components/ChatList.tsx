import React from "react";
import { List } from "antd";

interface Chat {
  id: string;
  name: string;
}

interface ChatListProps {
  chats: Chat[];
  onSelectChat: (chatId: string) => void;
}

const ChatList: React.FC<ChatListProps> = ({ chats, onSelectChat }) => {
  return (
    <List
      itemLayout="horizontal"
      dataSource={chats}
      renderItem={(chat) => (
        <List.Item onClick={() => onSelectChat(chat.id)} style={{ cursor: "pointer" }}>
          <List.Item.Meta title={chat.name} />
        </List.Item>
      )}
    />
  );
};

export default ChatList;
