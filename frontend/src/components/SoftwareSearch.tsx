import React, { useState } from "react";

const softwareList: string[] = [
  "AutoCAD",
  "Photoshop",
  "Illustrator",
  "Premiere Pro",
  "After Effects",
  "Visual Studio",
  "IntelliJ IDEA",
  "Eclipse",
  "NetBeans",
  "Unity",
  "Unreal Engine",
  "Blender",
  "Maya",
  "3ds Max"
];

export const SoftwareSearch: React.FC = () => {
  const [query, setQuery] = useState("");

  const filteredList = softwareList.filter(soft =>
    soft.toLowerCase().includes(query.toLowerCase())
  );

  return (
    <div style={{ fontFamily: "sans-serif", padding: "1rem", maxWidth: "400px" }}>
      <h2>Software Search</h2>

      <input
        type="text"
        placeholder="Type to search..."
        value={query}
        onChange={e => setQuery(e.target.value)}
        style={{
          width: "100%",
          padding: "0.5rem",
          marginBottom: "1rem",
          border: "1px solid #ccc",
          borderRadius: "6px"
        }}
      />

      <ul>
        {filteredList.length > 0 ? (
          filteredList.map((soft, index) => <li key={index}>{soft}</li>)
        ) : (
          <li>No software found</li>
        )}
      </ul>
    </div>
  );
};
