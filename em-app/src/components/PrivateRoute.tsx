import React from "react";
import { Navigate } from "react-router-dom";
import Cookies from "js-cookie";

// Интерфейс для PrivateRoute
interface PrivateRouteProps {
  element: React.ReactNode; // Компонент, который будет отрендерен
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ element }) => {
  const accessToken = Cookies.get("access_token");

  // Если токена нет, перенаправляем на страницу логина
  if (!accessToken) {
    return <Navigate to="/login" replace />;
  }

  // Если токен есть, рендерим переданный компонент
  return <>{element}</>;
};

export default PrivateRoute;
