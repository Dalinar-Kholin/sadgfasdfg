import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const useCheckCookie = () => {
    const navigate = useNavigate();

    useEffect(() => {
        const checkCookie = async () => {
            fetch('/api/checkCookie', {
                method: 'POST',
                headers: {
                    credentials: 'include'
                }
            }).then(response => {
                return response.json()
            }).then(
                data => {
                    if (!data.response) {
                        navigate('/login');
                    }
                }
            ).catch(error => {
                console.error('Error checking cookie:', error);
                navigate('/login');
            })
        };

        checkCookie().then(()=>
            console.log("logged")
        );
    }, [navigate]);
};

export default useCheckCookie;