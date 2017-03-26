/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Builds from './Builds';

document.body.className = 'root pt-dark';

export default class Main extends React.Component<void, void> {
	render(): JSX.Element {
		return <Builds/>;
	}
}
