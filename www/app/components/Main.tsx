/// <reference path="../References.d.ts"/>
import * as React from 'react';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import Styles from '../Styles';
import Alerts from './Alerts';
import Builds from './Builds';

document.body.style.backgroundColor = Styles.colors.background;

export default class Main extends React.Component<void, void> {
	render(): JSX.Element {
		return <MuiThemeProvider muiTheme={Styles.theme}>
			<div>
				<Builds/>
				<Alerts/>
			</div>
		</MuiThemeProvider>;
	}
}
