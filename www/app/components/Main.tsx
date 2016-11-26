/// <reference path="../References.d.ts"/>
import * as React from 'react';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import Styles from '../Styles';
import Builds from './Builds';

document.body.style.backgroundColor = Styles.colors.background;

export default class Main extends React.Component<void, void> {
	render(): JSX.Element {
		return <MuiThemeProvider muiTheme={Styles.theme}>
			<Builds/>
		</MuiThemeProvider>;
	}
}
