import React, { useState } from "react";

const ChatSSEPost = () => {
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState<{ role: string; content: string }[]>([]);
  const [waiting, setWaiting] = useState(false);

  const sendMessage = async () => {
    if (!input.trim()) return;

    setMessages((prev) => [...prev, { role: "user", content: input }]);
    setInput("");
    setWaiting(true);

    try {
      const response = await fetch("http://localhost:8080/prompt/sse", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ prompt: input }),
      });

      if (!response.ok || !response.body) {
        throw new Error("Invalid response");
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder("utf-8");

      let partial = "";
      let accumulated = "";

      while (true) {
        const { value, done } = await reader.read();
        if (done) break;

        partial += decoder.decode(value, { stream: true });

        // Split on SSE-style chunks (lines ending in `\n\n`)
        const events = partial.split(/\n\n/);
        partial = events.pop() || "";

        for (const event of events) {
          const line = event.trim().replace(/^data:\s*/, "");
          if (line === "[done]") {
            setMessages((prev) => [...prev, { role: "system", content: accumulated.trim() }]);
            setWaiting(false);
            return;
          }

          if (/^[.,!?;:]/.test(line)) {
            // If the line starts with punctuation, append directly
            accumulated += line;
          } else {
            // Otherwise, append with a space
            accumulated += (accumulated ? " " : "") + line;
          }
        }
      }
    } catch (err) {
      console.error("Error reading stream:", err);
      setWaiting(false);
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
        {waiting && (
          <div className="p-2 rounded-md max-w-xl bg-gray-100 self-start text-gray-500">
            Waiting for response...
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
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type your message..."
        />
        <button className="bg-blue-500 text-white px-4 py-1 rounded" type="submit">
          Send
        </button>
      </form>
    </div>
  );
};

export default ChatSSEPost;
