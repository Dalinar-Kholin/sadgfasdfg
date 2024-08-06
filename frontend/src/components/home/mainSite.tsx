import {ReactNode, useCallback, useEffect, useState} from "react";
import {
    Alert,
    AlertTitle,
    Button,
    CircularProgress,
    List,
    ListItemButton,
    ListItemText,
    TextField,
    Typography
} from "@mui/material";
import HurtResultForm from "../hurtResult/HurtResultForm.tsx";
import {hurtNames, IAllResult, IItemInstance, IItemToSearch} from "../../interfaces.ts";
import InputComp from "./inputField/InputField.tsx";
import stratchInputType from "./inputField/handleInputTypes/stracher.ts";
import Box from '@mui/material/Box';
import {getHurtResult, getMultipleHurtResult} from "./resultGrabbers.ts";

// niebieska obwódka -- najtańszy pakiet
// zielona obwódka -- najtanicej za produkt
// czerwona obwódka -- najdrożej za produkt


export default function MainSite() {
    // region zmienne
    const [Ean, setEan] = useState<string>("")
    let availableHurt = localStorage.getItem("availableHurt") ? parseInt(localStorage.getItem("availableHurt") || "0") : 15

    const [componentHashTable, setComponentHashTable] = useState<Map<hurtNames, ReactNode>>(new Map<hurtNames, ReactNode>())


    const [errorMessage, setErrorMessage] = useState<string>("")

    const [lowerHurt, setLowerHurt] = useState<hurtNames>(hurtNames.none)
    const [isLoadingProduct, setIsLoadingProduct] = useState<boolean>(false)


    const [prodToSearch, setProdToSearch] = useState<IItemToSearch[]>([])


    const [optItems, setOptItems] = useState<IItemInstance[]>([])
    const [allResult, setAllResult] = useState<IAllResult[]>([])


    const [fileName, setFileName] = useState<string>("")

    // endregion

    // region pozwala na przeciąganie plików
    const onDrop = useCallback((event: DragEvent) => {
        event.preventDefault();
        const file = event.dataTransfer?.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e: ProgressEvent<FileReader>) => {
                const result = e.target?.result;
                if (result && typeof result === 'string') {
                    setFileName(file.name);
                    setProdToSearch(stratchInputType(result, file.name)(result));
                }
            };

            reader.readAsText(file);
        }
    }, []);

    const onDragOver = useCallback((event: DragEvent) => {
        event.preventDefault();
    }, []);

    useEffect(() => {

        window.addEventListener('dragover', onDragOver);
        window.addEventListener('drop', onDrop);

        return () => {
            window.removeEventListener('dragover', onDragOver);
            window.removeEventListener('drop', onDrop);
        };
    }, [onDrop, onDragOver]);
    // endregion

    const changeResultComp = (ean: string) => {
        const newComponentHashTable = new Map<hurtNames, ReactNode>()

        if (allResult.filter((item) => item.ean === ean).length === 0) {
            newComponentHashTable.set(hurtNames.none,
                <HurtResultForm
                    name={hurtNames[hurtNames.none]}
                    priceForPack={-1}
                    princeForOne={-1}
                    productsInPack={-1}/>
            )
        }

        allResult.filter((item) => item.ean === ean).map((newItem) => {

            newItem.result.filter(i => i.Item.priceForOne !== -1).map((newItem) => {
                newComponentHashTable.set(newItem.hurtName,
                    <HurtResultForm
                        name={hurtNames[newItem.hurtName]}
                        priceForPack={newItem.Item.priceForPack}
                        princeForOne={newItem.Item.priceForOne}
                        productsInPack={newItem.Item.productsInPack}/>
                )
            })
        })
        setComponentHashTable(newComponentHashTable)
    }


    useEffect(() => { // ładowanie mapy komponentów
        const newHashTable = new Map<hurtNames, ReactNode>();

        for (let i = 1; i <= availableHurt; i <<= 1) {
            if ((availableHurt & i) !== 0) {
                newHashTable.set(i, (
                    <HurtResultForm
                        name={hurtNames[i]}
                        priceForPack={-1}
                        princeForOne={-1}
                        productsInPack={-1}
                    />
                ));
            }
        }
        setComponentHashTable(newHashTable);
    }, []) // tworzenie nowej hashMapy aby nie świeciło pustkami

    useEffect(() => {
        setIsLoadingProduct(true)
        try {
            getMultipleHurtResult(prodToSearch).then(data => {
                const newOptItems: IItemInstance[] = []
                const newAllResult: IAllResult[] = []
                prodToSearch.map((item) => {
                    const ItemsMatchEan = data.get(item.Ean)
                    if (ItemsMatchEan) {
                        const newMin = ItemsMatchEan.filter((element) => {
                            return element.Item.priceForOne !== -1
                        })
                        if (!newMin || newMin.length === 0) {
                            return;
                        }

                        const xd = newMin.reduce((prev, current) => {
                            return prev.Item.priceForOne < current.Item.priceForOne ? prev : current
                        })
                        newOptItems.push({
                            name: item.Name,
                            ean: item.Ean,
                            item: xd.Item,
                            count: item.Amount,
                        })
                        newAllResult.push({
                            ean: item.Ean,
                            result: ItemsMatchEan
                        })
                    } else {
                        newOptItems.push({
                            name: item.Name,
                            ean: item.Ean,
                            item: {
                                hurtName: hurtNames.none,
                                priceForPack: -1,
                                priceForOne: -1,
                                productsInPack: -1,
                            },
                            count: item.Amount,
                        })
                    }
                })
                setOptItems(newOptItems)
                setAllResult(newAllResult)
                setIsLoadingProduct(false)
            })
        } catch (e: any){
            setErrorMessage(e.message)
            setIsLoadingProduct(false)
        }
        // zapisanie ich w optItems
    }, [prodToSearch])


    return (
        <>
            <h1>skanuj pojedyńczo</h1>
            <p></p>

            <TextField autoComplete={"off"} id="filled" label="Kod Ean" placeholder="Ean" value={Ean}
                       onChange={e => setEan(e.target.value)} onKeyDown={e => {
                if (e.key === "Enter") {
                    setIsLoadingProduct(true)
                    try {
                        getHurtResult(Ean).then(data => {
                            const newMap = new Map<hurtNames, ReactNode>()
                            let i = 0
                            data.map((item) => {
                                if (item.priceForOne !== -1) {
                                    i += 1
                                    newMap.set(item.hurtName, (
                                        <HurtResultForm
                                            name={hurtNames[item.hurtName]}
                                            priceForPack={item.priceForPack}
                                            princeForOne={item.priceForOne}
                                            productsInPack={item.productsInPack}
                                        />
                                    ))
                                }
                            })
                            if (i === 0) {
                                newMap.set(hurtNames.none, (
                                    <HurtResultForm
                                        name={hurtNames[hurtNames.none]}
                                        priceForPack={-1}
                                        princeForOne={-1}
                                        productsInPack={-1}
                                    />
                                ))
                            } else {
                                setComponentHashTable(newMap)
                            }

                            setLowerHurt(Math.max(...data.map(item => item.priceForOne)) === -1 ?
                                hurtNames.none : data.find(item => item.priceForOne === Math.max(...data.map(item => item.priceForOne)))?.hurtName || hurtNames.none)
                            setIsLoadingProduct(false)
                        });
                    }catch (e: any){
                        setErrorMessage(e.message)
                        setIsLoadingProduct(false)
                    }
                }}}
            />

            <p></p>
            {hurtNames[lowerHurt]}
            <Box style={{display: "flex", justifyContent: "space-around", alignItems: "center"}}>
                {!isLoadingProduct ? (
                    <div className={"hurtResults"}
                         style={{margin: "10px", padding: "10px", display: "grid", gap: "10px"}}>
                        {
                            Array.from(componentHashTable.values()).map((element) => {
                                return element
                            })

                        }
                    </div>
                ) : <Box sx={{display: 'flex', padding: "20px"}}>
                    <CircularProgress/>
                </Box>
                }
                <List component="nav"
                      sx={{width: "20%", overflow: "scroll", scrollbarWidth: "none", maxHeight: "600px"}}>
                    {prodToSearch.map((item) => {
                        return (
                            <ListItemButton onClick={() => {
                                changeResultComp(item.Ean)
                            }}>
                                <ListItemText primary={item.Name}/>
                            </ListItemButton>
                        )
                    })}
                </List>

            </Box>
            <Button sx={{margin: "20px", padding: "5px"}}
                    variant="outlined" color="error" onClick={() => {
                setProdToSearch([])
            }}>Wyczyść Listę</Button>
            {optItems.length !== 0 ?
                <Button variant="contained" color="success" sx={{margin: "20px", padding: "5px"}} onClick={() => {
                    setErrorMessage("")
                    fetch("/api/makeOrder", {
                        method: "POST",
                        body: JSON.stringify({Items: optItems.map(item => {
                            return {
                                Ean: item.ean,
                                Amount: item.count,
                                HurtName: item.item.hurtName
                            }
                            })}),
                        credentials: "include",
                        headers: {
                            "Content-Type": "application/json",
                        }

                    }).then(response => {
                        if (response.status !== 200) {
                            throw new Error("nie udało się złożyć zamówienia")
                        }
                        return response.json()
                    }).then(data => {
                        console.log(data)
                    }).catch(err => {
                        setErrorMessage(err)
                        throw new Error(err);
                    })
                }}>
                    Rozdziel produkty do koszyka
                </Button> : null}
            {errorMessage !== "" ?
                <Alert severity="error">
                    <AlertTitle>Error</AlertTitle>
                    {errorMessage}
                </Alert>
                : null}

            {fileName && <Box marginTop="1rem" width="100%">
                <Typography variant="h6" component="h2" gutterBottom>
                    {"przetwarzany plik := " + fileName}
                </Typography>
            </Box>}

            <InputComp setItem={prod => setProdToSearch(prod)} setName={name => setFileName(name)}/>

        </>
    )
}



