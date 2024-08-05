import {useState} from "react";
import {Alert, AlertTitle, Button, TextField} from "@mui/material";
import "./login.css"
import {useNavigate} from "react-router-dom";

function getCookie(name: string) {
    let cookieArr = document.cookie.split(";");

    for (let i = 0; i < cookieArr.length; i++) {
        let cookiePair = cookieArr[i].split("=");

        // Usunięcie spacji na początku nazwy ciasteczka
        let cookieName = cookiePair[0].trim();
        if (cookieName === name) {
            return cookiePair[1];
        }
    }
    return null;
}

export default function LoginForm(){
    if (getCookie("accessToken") !== null){
        //return <h1>already logged in</h1>
    }

    const [username, setUsername] = useState<string>();
    const [password, setPassword] = useState<string>();
    const [isProperData, setIsProperData] = useState<boolean>(true)
    const [errorMessage, setErrorMessage] = useState<string>("")

    const navigate = useNavigate()

    return (
        <>
            <h1>Login Form</h1>
            <form onSubmit={e => {
                e.preventDefault()
                // logowanie sie
                const LoginData = {
                    username : username,
                    password : password,
                }
                fetch('/api/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(LoginData),
                }).then(response => {
                    if (!response.ok) {
                        // Próbujemy pobrać dane błędu jako JSON
                        return response.json().then(errorData => {
                            setErrorMessage(errorData.error || "błąd sieci")
                            setIsProperData(false)
                            throw new Error('Error: ' + errorData.error);
                        });
                    }
                    return response.json();
                }).then(data => {
                    setIsProperData(true)
                    setUsername("")
                    setPassword("")
                    // dodać ciasteczka
                    // przenieść na stronę główną
                    navigate("/main")
                    console.log(data)
                }).catch(error => {
                    console.error('There has been a problem with your fetch operation:', error);
                })

            }}>
                <TextField
                    autoComplete={"off"}
                    id="filled"
                    label="Username"
                    placeholder="username"
                    value={username}
                    onChange={e => setUsername(e.target.value)}
                />
                <p></p>
                <TextField
                    id="outlined-password-input"
                    label="Password"
                    type="password"
                    autoComplete="current-password"
                    value={ password }
                    onChange={e => setPassword(e.target.value)}
                />
                <p></p>
                {isProperData ? <div></div> : <Alert severity="error">
                    <AlertTitle>Error</AlertTitle>
                    {errorMessage}
                </Alert>}
                <p></p>
                <Button variant="contained" type={"submit"}>
                    Login
                </Button>
            </form>
        </>
    )

}