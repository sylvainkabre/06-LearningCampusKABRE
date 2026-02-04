import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
    {/* main div */}
      <div className="flex">
        {/* sidebar */}
        <div className="w-64 bg-gray-200 p-4">
          Sidebar
        </div>
        {/* body */}
        <div className="flex-1 bg-gray-100 p-4">
        Body
        </div>
      </div>
    </>
  )
}

export default App
