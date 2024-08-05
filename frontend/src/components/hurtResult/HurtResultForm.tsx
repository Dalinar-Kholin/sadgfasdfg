import {TextField} from "@mui/material";
import "./hurtResult.css"

interface IHurtResultForm {
    name: string,
    princeForOne: number,
    priceForPack: number,
    productsInPack: number
}

export default function HurtResultForm({name, princeForOne, priceForPack, productsInPack}: IHurtResultForm) {
    return (
        <>
            <div className={"masterHurtDiv"} >
                <div className={"hurtResultDiv"}>
                    <TextField
                        disabled
                        id="outlined-disabled"
                        label="nazwa Hurtowni"
                        value={name}
                    />
                    <TextField
                        disabled
                        id="outlined-disabled"
                        label="cena za sztukę"
                        value={princeForOne === -1 ? "Brak produktu" : princeForOne}
                    />
                    <TextField
                        disabled
                        id="outlined-disabled"
                        label={"Cena za " + (productsInPack === -1 ? "" : productsInPack)}
                        value={priceForPack === -1 ? "Brak produktu" : priceForPack.toFixed(2)}
                    />
                </div>
            </div>
        </>
    )
}