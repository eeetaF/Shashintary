import os
import requests
from flask import Flask, request, jsonify
from dotenv import load_dotenv

load_dotenv()
API_KEY = os.getenv("OPENROUTER_API_KEY")

app = Flask(__name__)

@app.route('/generate', methods=['POST'])
def generate():
    data = request.get_json()

    prompt = data.get("prompt")
    min_tokens = data.get("min_tokens", 200)
    max_tokens = data.get("max_tokens", 400)
    temperature = data.get("temperature", 0.7)

    if not prompt:
        return jsonify({"error": "No prompt provided"}), 400

    headers = {
        "Authorization": f"Bearer {API_KEY}",
        "Content-Type": "application/json",
        "HTTP-Referer": "http://localhost",
    }

    payload = {
        "model": "openai/gpt-3.5-turbo",
        "messages": [
            {"role": "system", "content": "You are a chess commentator. Use the template provided to generate a comment. Answer should be 150-300 symbols long. Use your creativity."},
            {"role": "user", "content": prompt}
        ],
        "temperature": temperature,
        "max_tokens": max_tokens
    }

    try:
        response = requests.post("https://openrouter.ai/api/v1/chat/completions", headers=headers, json=payload)
        response.raise_for_status()
        print("RESPONSE JSON:", response.json())  # <-- debug output
        result = response.json()["choices"][0]["message"]["content"]
        return jsonify({"result": result})
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    app.run(host="localhost", port=53004)
