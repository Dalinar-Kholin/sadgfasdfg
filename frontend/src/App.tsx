import LoginForm from "./components/login/loginForm.tsx";
import {BrowserRouter as Router, Navigate, Route, Routes} from 'react-router-dom';
import MainSite from "./components/home/mainSite.tsx";
import AppBarCustomed from "./components/appBar/appBar.tsx";
import {createTheme, CssBaseline, ThemeProvider} from "@mui/material";
import useCheckCookie from "./customHook/useCheckCookie.ts";


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
     // sporawdza czy jest G a jeżeli nie to przekierowuje na stronę logowania

    /*useEffect(() => { // powino być odpowiedzialne za wysłanie requesta do serwera z informacją o opuszczeniu strony -- wylogowanie i brak potrzeby sprawdzania sesji
        const handleBeforeUnload = () => {
            // Użycie navigator.sendBeacon do wysłania danych
            fetch('https://127.0.0.1:8080/api/exit', {
                method: "POST",
                credentials: "include"
            })
        };

        window.addEventListener('beforeunload', handleBeforeUnload);

        // Usuwanie zdarzenia po odmontowaniu komponentu
        return () => {
            window.removeEventListener('beforeunload', handleBeforeUnload);
        };
    }, []);*/ // jeden wielki chuj, nie ma sensu usuwać ciasteczek, ponieważ użytkownik przy każdym realodzie musiłby restartować sesję i logować się ponownie

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
                        <Route path={"/*"} element={<Navigate to={"/main"}/>}/>
                    </Routes>
                </Router>
            </ThemeProvider>
        </div>
    )
}

export default App
