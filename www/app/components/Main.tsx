/// <reference path="../References.d.ts"/>
import * as React from 'react';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import darkBaseTheme from 'material-ui/styles/baseThemes/darkBaseTheme';
import Builds from './Builds';

export default class Main extends React.Component<void, void> {
	render(): JSX.Element {
		return <MuiThemeProvider muiTheme={getMuiTheme(darkBaseTheme)}>
			<Builds/>
		</MuiThemeProvider>;
	}
}
