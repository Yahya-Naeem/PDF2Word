"use client"; // Needed for Next.js App Router (Client Component)

import { useState } from "react";

const UploadForm = () => {
  const [file, setFile] = useState<File | null>(null);
  const [downloadUrl, setDownloadUrl] = useState("");

  const handleUpload = async () => {
    if (!file) return alert("Please select a PDF file");

    const formData = new FormData();
    formData.append("pdf", file);

    const response = await fetch("http://localhost:8080/upload", {
      method: "POST",
      body: formData,
    });

    if (response.ok) {
      const fileUrl = await response.text();
      setDownloadUrl(fileUrl);
    } else {
      alert("Upload failed!");
    }
  };

  return (
    <div className="p-6 bg-white shadow-lg rounded-lg text-center w-96">
      <h2 className="text-xl font-bold mb-4">Upload a PDF</h2>
      <input
        type="file"
        accept="application/pdf"
        onChange={(e) => setFile(e.target.files?.[0] || null)}
        className="mb-4 border p-2 w-full"
      />
      <button
        onClick={handleUpload}
        className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 w-full"
      >
        Convert
      </button>

      {downloadUrl && (
        <div className="mt-4">
          <a href={downloadUrl} className="text-green-600 underline">
            Download Word File
          </a>
        </div>
      )}
    </div>
  );
};

export default UploadForm;
