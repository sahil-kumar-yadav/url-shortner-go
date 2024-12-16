"use client"

import { useState } from "react";
import axios from 'axios';

export default function Home() {

  const [originalURL, setOriginalURL] = useState('');
  const [shortenedURL, setShortenedURL] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setShortenedURL('');

    try {
      const response = await axios.post('http://localhost:3000/shorten', {
        url: originalURL,
      });
      setShortenedURL(response.data.shortedURL);
    } catch (err) {
      setError('Failed to shorten the URL. Please try again.',err);
      console.log(err)
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen">
      <h1 className="text-4xl font-bold mb-6">URL Shortener</h1>
      <form
        onSubmit={handleSubmit}
        className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
      >
        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-bold mb-2">
            Original URL
          </label>
          <input
            type="url"
            placeholder="Enter the URL"
            value={originalURL}
            onChange={(e) => setOriginalURL(e.target.value)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
            required
          />
        </div>
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          Shorten
        </button>
      </form>
      {shortenedURL && (
        <p className="text-green-500">
          Shortened URL:{" "}
          <a
            href={shortenedURL}
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-700 underline"
          >
            {shortenedURL}
          </a>
        </p>
      )}
      {error && <p className="text-red-500 mt-4">{error}</p>}
    </div>
  );
}
