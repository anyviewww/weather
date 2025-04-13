import React, { useState } from 'react';

function App() {
  const [city, setCity] = useState('');
  const [data, setData] = useState(null);

  const fetchData = async () => {
    try {
        const response = await fetch(`http://localhost:8080/weather/${city}`);
        if (!response.ok) {
            throw new Error(`Error: ${response.statusText}`);
        }
        const result = await response.json();
        setData(result);
    } catch (error) {
        console.error('Error fetching data:', error);
    }
};


  return (
    <div style={{ padding: '20px' }}>
      <h1>Weather and Time</h1>
      <input
        type="text"
        placeholder="Enter city name"
        value={city}
        onChange={(e) => setCity(e.target.value)}
      />
      <button onClick={fetchData}>Get Weather</button>

      {data && (
        <div style={{ marginTop: '20px' }}>
          <h2>{data.city}</h2>
          <p>Temperature: {data.weather}°C</p>
          <p>Time: {data.time}</p>
        </div>
      )}
    </div>
  );
}

export default App;
