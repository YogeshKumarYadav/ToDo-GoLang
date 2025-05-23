import { Stack } from "@chakra-ui/react"
import Navbar from "./components/Navbar"
import TodoForm from "./components/TodoForm"
import TodoList from "./components/TodoList"

export const BASE_URL = import.meta.env.MODE === "development" ? "http://localhost:5000/api" : "/api";

function App() {
  return (
    <Stack height={"100vh"}>
      <Navbar/>
      <Stack style={{width: "60vh", minWidth: "20vh", margin: "0 auto"}}>
        <TodoForm/>
        <TodoList/>
      </Stack>
    </Stack>
  )
}

export default App
