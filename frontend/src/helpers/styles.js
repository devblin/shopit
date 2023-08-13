import { makeStyles } from '@mui/styles';
import { createTheme, InputBase } from '@mui/material';
import { styled, alpha } from '@mui/material/styles';
import { PRIMARY_COLOR } from './constants';

const theme = createTheme({
	palette: {
		primary: {
			main: PRIMARY_COLOR,
		}
	}
});

const generalStyles = makeStyles({
	container: {
		display: "flex !important",
		alignItems: "center !important",
		minHeight: "750px !important",
		flexDirection: "column",
		flexWrap: "wrap",
		marginTop: "100px"
	},
});

const itemStyles = makeStyles({
	image: {
		maxHeight: "400px",
		height: "100%",
		objectFit: "contain !important"
	},
	cardContent: {
		padding: "15px !important"
	},
	card: {
		margin: 1
	}
})

const navbarStyles = {
	searchArea: styled('div')(({ theme }) => ({
		position: 'relative',
		borderRadius: theme.shape.borderRadius,
		backgroundColor: alpha(theme.palette.common.white, 0.15),
		'&:hover': {
			backgroundColor: alpha(theme.palette.common.white, 0.25),
		},
		marginRight: theme.spacing(2),
		marginLeft: theme.spacing(3),
		width: '100%'
	})),
	searchIconWrapper: styled('div')(({ theme }) => ({
		padding: theme.spacing(0, 2),
		height: '100%',
		position: 'absolute',
		pointerEvents: 'none',
		display: 'flex',
		alignItems: 'center',
		justifyContent: 'center',
	})),
	styledInputBase: styled(InputBase)(({ theme }) => ({
		color: 'inherit',
		'& .MuiInputBase-input': {
			padding: theme.spacing(1, 1, 1, 0),
			// vertical padding + font size from searchIcon
			paddingLeft: `calc(1em + ${theme.spacing(4)})`,
			transition: theme.transitions.create('width'),
			width: '100%',
			[theme.breakpoints.up('md')]: {
				width: '20ch',
			},
		},
	}))
}

export {
	theme,
	generalStyles,
	itemStyles,
	navbarStyles
};