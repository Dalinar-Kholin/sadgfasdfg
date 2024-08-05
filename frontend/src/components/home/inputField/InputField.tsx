import {ChangeEvent} from "react";
import {Button, Container,  Typography} from "@mui/material";
import Box from "@mui/material/Box";
import {IItemToSearch} from "../../../interfaces.ts";
import stratchInputType from "./handleInputTypes/stracher.ts";


interface IInputField{
    setItem : (items : IItemToSearch[]) => void // zjąć się handlerką danych w inputField
    setName : (name : string) => void
}


export default function InputComp({setItem,setName}: IInputField){

    const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e: ProgressEvent<FileReader>) => {
                const result = e.target?.result;
                if (result && typeof result === 'string') {
                    setName(file.name);
                    setItem(stratchInputType(result, file.name)(result))
                }
            };

            reader.readAsText(file);


        }
    };

    return (
        <Container maxWidth="md" style={{marginTop: '2rem'}}>
            <Box display="flex" flexDirection="column" alignItems="center">
                <Typography variant="h4" component="h1" gutterBottom>
                    File Reader
                </Typography>
                <input
                    accept={".txt" || ".pdf" || ".xlsx"}
                    style={{display: 'none'}}
                    id="file-input"
                    type="file"
                    onChange={handleFileChange}
                />
                <label htmlFor="file-input">
                    <Button variant="contained" color="primary" component="span">
                        Upload File
                    </Button>
                </label>
            </Box>
        </Container>
    );
}