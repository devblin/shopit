import { Container } from "@mui/material";
import { useEffect } from "react";
import { Route, Switch } from "wouter"
import { generalStyles } from "../helpers/styles";
import Items from "./pages/Items";
import AddItem from "./pages/items/AddItem";
import EditItem from "./pages/items/EditItem";
import ViewItem from "./pages/items/ViewItem";
import NotFound from "./pages/NotFound";

function Main() {
    const generalClasses = generalStyles()

    useEffect(() => {
        document.title = "Shopit"
    });

    return (
        <Container className={generalClasses.container}>
            <Switch>
                <Route key={1} path="/add-item" component={AddItem}></Route>
                <Route key={2} path="/view-item/:itemId" component={ViewItem}></Route>
                <Route key={3} path="/edit-item/:itemId" component={EditItem}></Route>
                <Route key={4} path="/" component={Items}></Route>
                <Route key={5} component={NotFound}></Route>
            </Switch>
        </Container >
    );
}

export default Main;