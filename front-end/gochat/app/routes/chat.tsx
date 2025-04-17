import { useState } from "react";

export default function ChatPage() {
  const [messages, setMessages] = useState([
    { role: "assistant", content: "Hello! How can I help you today?" },
  ]);
  const [input, setInput] = useState("");

  const sendMessage = async () => {
    if (!input.trim()) return;

    const userMessage = { role: "user", content: input };
    const newMessages = [...messages, userMessage];
    setMessages(newMessages);
    setInput("");

    try {
      const response = await fetch("http://localhost:8080/query", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: `
            query Prompt($prompt: String!) {
                promptResponse(prompt: $prompt)
            }
          `,
          variables: {
            prompt: input, // Send the user input as the prompt
          },
        }),
      });

      const { data, errors } = await response.json();

      if (errors) {
        console.error("GraphQL errors", errors);
        return;
      }

      const reply = data.promptResponse;
      setMessages((prev) => [...prev, { role: "assistant", content: reply }]);
    } catch (err) {
      console.error("Network error", err);
    }
  };

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
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type your message..."
        />
        <button className="bg-blue-500 text-white px-4 py-1 rounded" type="submit">
          Send
        </button>
      </form>
    </div>
  );
}
