/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as MuiStyles from 'material-ui/styles';
import Styles from '../Styles';
import Builds from './Builds';

document.body.style.backgroundColor = Styles.colors.background;

export default class Main extends React.Component<void, void> {
	render(): JSX.Element {
		return <MuiStyles.MuiThemeProvider muiTheme={Styles.theme}>
			<Builds/>
		</MuiStyles.MuiThemeProvider>;
	}
}
