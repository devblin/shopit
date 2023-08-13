import { Card, CardHeader, CardMedia, CardContent, Typography, Stack, Button, CardActions } from '@mui/material';
import { itemStyles } from "../../../helpers/styles"
import { getImageBySize, getPriceComponent, getStockComponent } from '../../../helpers/get';
import { SEVERITY } from "../../../helpers/constants"
import { Edit } from '@mui/icons-material';

export default function Item(props) {
    const classes = itemStyles();

    return (
        <Card className={classes.card}>
            <CardHeader
                title={props.name}
            />
            <CardMedia
                className={classes.image}
                component="img"
                image={getImageBySize(props.image, 400)}
                alt="item-image"
            />
            <CardContent className={classes.cardContent}>
                <Typography fontSize={20} color="text.secondary">
                    <b>Category:</b> {props.category}
                </Typography>
                <Stack>
                    {getStockComponent(props.stock, 20)}
                    {getPriceComponent(props.price)}
                </Stack>
            </CardContent>
            <CardActions>
                <Button
                    fullWidth
                    variant='contained'
                    href={`/edit-item/${props.id}`}
                    color={SEVERITY.WARN}
                    size="large"
                >
                    <Edit />
                    <code>Edit</code>
                </Button>
            </CardActions>
        </Card>
    );
}