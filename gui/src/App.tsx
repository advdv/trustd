import { useState } from "react";
import trustdLogo from "/trustd.svg";
import "./App.css";

function App() {
  const [count, setCount] = useState(0);

  return (
    <>
      <div>
        <img src={trustdLogo} className="logo" alt="Vite logo" />
      </div>
      <h1 className="text-3xl font-bold underline">Trustd</h1>
      <div className="card">
        <button
          type="button"
          onClick={() => {
            setCount((count) => count + 1);
          }}
        >
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  );
}

export default App;
