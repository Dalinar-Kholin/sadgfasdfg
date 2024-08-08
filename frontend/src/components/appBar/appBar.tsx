import {AppBar, Button, Container, Toolbar} from "@mui/material";
import Box from "@mui/material/Box";
import {useNavigate} from "react-router-dom";
import "./appBar.css"
import {PATH} from "../../interfaces.ts";



/*interface IAppBarCustomed{
    iconLink : string
}*/

export default function AppBarCustomed(/*{iconLink}:IAppBarCustomed*/){

    const navigate = useNavigate()

    const Style={
        logoIcon:{
            borderRadius: '50%',
            height:'50px',
            width:'50px',
        }

    }


    const pages : PATH[] = [ "ustawienia", "cennik"]
    return(
        <>
            <AppBar position="sticky" id={"appBarComp"}>
                <Container maxWidth="xl">
                    <Toolbar disableGutters>
                        <Box className="photo" component="img" src={"./assets/nicea.jpeg"} alt="Logo" style={Style.logoIcon} />
                        <Box >
                            {pages.map((page) => (
                                <Button
                                    color="inherit"
                                    key={page}
                                    onClick={()=>{
                                        navigate('/' + page);
                                    }}
                                >
                                    {page}
                                </Button>
                            ))}
                        </Box>
                    </Toolbar>
                </Container>
            </AppBar>
        </>
    )
}