import {IItemToSearch} from "../../../../interfaces.ts";

export default function handleKc(input : string) : IItemToSearch[]{
    return input.split("\n").map((line) => {
        const row = line.split(';')
        if (row.length <= 2) {
            return {Ean: "0", Amount: 0, Name: ""}
        }

        return {Ean: row[0], Amount: +row[1], Name: row[3]} // row[2] to cena
    }).filter(i => i.Ean !== "0" && i.Amount !== 0)


}