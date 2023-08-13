import { useEffect } from "react";
import { theme } from "./helpers/styles";
import { ThemeProvider } from "@mui/material/styles";
import NavBar from "./components/custom/NavBar";
import Main from "./components/Main";
import FootBar from "./components/custom/FootBar";

function App() {
	useEffect(() => {
		document.title = process.env.REACT_APP_NAME
	}, []);

	return (
		<ThemeProvider theme={theme}>
			<NavBar></NavBar>
			<Main></Main>
			<FootBar></FootBar>
		</ThemeProvider>
	);
}

export default App;