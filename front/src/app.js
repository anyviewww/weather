import React, { useEffect, useState } from 'react';

function App() {
  const [data, setData] = useState({ time: "Loading...", city: "", temp: "" });

  useEffect(() => {
    fetch('http://localhost:8080')
      .then(response => response.json())
      .then(data => setData(data))
      .catch(error => setData({ time: "Error fetching data", city: "", temp: "" }));
  }, []);

  return (
    <div>
      <h1>Current Time and Weather</h1>
      <p>Time: {data.time}</p>
      <p>City: {data.city}</p>
      <p>Temperature: {data.temp} °C</p>
    </div>
  );
}

export default App;