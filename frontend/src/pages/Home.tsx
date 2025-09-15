import { useState, type FC } from "react";

const Home: FC = () => {
    const [count, setCount] = useState(0)

  return (
    <section>
      <h2>Home</h2>
      <p>Contador: {count}</p>
      <button onClick={() => setCount(c => c + 1)}>Sumar</button>
    </section>
  )
}

export default Home
