import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
    index("routes/home.tsx"),
    route("chat", "routes/chat.tsx"),
    route("chatws", "routes/chatws.tsx"),
    route("chatsse", "routes/chatsse.tsx"),
] satisfies RouteConfig;
