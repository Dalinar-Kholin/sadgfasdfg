import LoginForm from "./components/login/loginForm.tsx";
import {BrowserRouter as Router, Navigate, Route, Routes} from 'react-router-dom';
import MainSite from "./components/home/mainSite.tsx";
import AppBarCustomed from "./components/appBar/appBar.tsx";
import {createTheme, CssBaseline, ThemeProvider} from "@mui/material";
import useCheckCookie from "./customHook/useCheckCookie.ts";
import Tariff from "./components/tariff/Tariff.tsx";


const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});


function CheckCookie() {
    useCheckCookie();
    return <></>
}

function App() {

    return (
        <div>
            <ThemeProvider theme={darkTheme}>
                <CssBaseline/>
                <Router>
                    <CheckCookie/>
                    {location.pathname === "/login" ? <></> : <AppBarCustomed />}
                    <Routes>
                        <Route path={"/login"} element={<LoginForm/>}/>
                        <Route path={"/main"} element={<MainSite/>}/>
                        <Route path={"/cennik"} element={<Tariff/>}/>
                        <Route path={"/ustawienia"} element={<MainSite/>}/>
                        <Route path={"/*"} element={<Navigate to={"/main"}/>}/>
                    </Routes>
                </Router>
            </ThemeProvider>
        </div>
    )
}

export default App
