import React, { useState, useEffect, useRef } from "react";

const ChatApp = () => {
  const socketRef = useRef<WebSocket>(null); // Ref to keep WebSocket instance persistent
  const [input, setInput] = useState(""); // User input message

  type Message = {
    role: "user" | "assistant" | "system";
    content: string;
  };

  const [messages, setMessages] = useState<Message[]>([]); // Store messages (both user and LLM responses)
  const [waiting, setWaiting] = useState(false); // Flag to indicate if the server is responding

  const updateMessages = (prevMessages: Message[], event: { data: string }): Message[] => {
    const newMessages = [...prevMessages];
    const lastMessage = newMessages[newMessages.length - 1];
  
    if (lastMessage && lastMessage.role === "system") {
      if (!lastMessage.content.includes(event.data)) {
        lastMessage.content += event.data;
      }
    } else {
      newMessages.push({ role: "system", content: event.data });
    }
  
    return newMessages;
  };
  
  // Initialize WebSocket connection once, and keep it throughout the component's lifecycle
  const initializeWebSocket = () => {
    if (socketRef.current) return; // Avoid re-initializing if it's already done

    socketRef.current = new WebSocket("ws://localhost:8080/prompt/ws");

    socketRef.current.onopen = () => {
      console.log("WebSocket connection established");
    };

    socketRef.current.onmessage = (event) => {
      console.log("Received from server:", event.data);

      // If the response is not "[done]", it means the server is still sending data
      if (event.data !== "[done]") {
        setMessages((prevMessages) => updateMessages(prevMessages, event));
      } else {
        // Mark the waiting flag as false when the server signals it's done
        console.log("Server indicates the response is done.");
        setWaiting(false);
      }
    };

    socketRef.current.onerror = (error) => {
      console.log("WebSocket error:", error);
    };

    socketRef.current.onclose = () => {
      console.log("WebSocket connection closed");
      // Clean up by resetting the socket reference
      socketRef.current = null;
    };
  };

  // Handle message input change
  const handleInputChange = (e) => {
    setInput(e.target.value);
  };

  // Handle message submission to send to the WebSocket server
  const sendMessage = () => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN && input.trim()) {
      setMessages((prevMessages) => [
        ...prevMessages,
        { role: "user", content: input },
      ]);

      setWaiting(true);

      socketRef.current.send(input);
      setInput("");
    } else {
      console.log("WebSocket is not open or input is empty.");
    }
  };

  // Clean up the WebSocket connection on component unmount or before unload
  useEffect(() => {
    initializeWebSocket();

    const handleBeforeUnload = () => {
      console.log("Window is unloading, closing WebSocket.");
      if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
        socketRef.current.close();
      }
    };

    window.addEventListener("beforeunload", handleBeforeUnload);

    // Cleanup on component unmount
    return () => {
      window.removeEventListener("beforeunload", handleBeforeUnload);
      if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
        socketRef.current.close();
      }
    };
  }, []); // Only run once on mount/unmount

  return (
    <div className="flex flex-col h-screen">
      <div className="flex-1 overflow-y-auto p-4 space-y-2 bg-gray-50">
        {messages.map((msg, idx) => (
          <div
            key={idx}
            className={`p-2 rounded-md max-w-xl ${
              msg.role === "user" ? "bg-blue-200 self-end" : "bg-gray-200 self-start"
            }`}
          >
            {msg.content}
          </div>
        ))}

        {waiting && (
          <div className="p-2 rounded-md max-w-xl bg-gray-100 self-start">
            <span className="text-gray-500">Waiting for response...</span>
          </div>
        )}
      </div>

      <form
        className="flex p-4 border-t bg-white"
        onSubmit={(e) => {
          e.preventDefault();
          sendMessage();
        }}
      >
        <input
          className="flex-1 border rounded px-2 py-1 mr-2"
          value={input}
          onChange={handleInputChange}
          placeholder="Type your message..."
        />
        <button className="bg-blue-500 text-white px-4 py-1 rounded" type="submit">
          Send
        </button>
      </form>
    </div>
  );
};

export default ChatApp;
