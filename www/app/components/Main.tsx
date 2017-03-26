/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as BuildActions from '../actions/BuildActions';
import Loading from './Loading';
import Builds from './Builds';
import Nodes from './Nodes';

document.body.className = 'root pt-dark';

interface State {
	page: string;
}

export default class Main extends React.Component<void, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			page: 'builds',
		};
	}

	sync = (): void => {
		switch (this.state.page) {
			case 'builds':
				BuildActions.sync();
				break;
		}
	}

	render(): JSX.Element {
		let page: JSX.Element;

		switch (this.state.page) {
			case 'builds':
				page = <Builds/>;
				break;
			case 'nodes':
				page = <Nodes/>;
				break;
		}

		return <div>
			<nav className="pt-navbar">
				<div className="pt-navbar-group pt-align-left">
					<div className="pt-navbar-heading">AutoABS</div>
					<Loading size="small"/>
					<div className="pt-navbar-heading">Builds</div>
					<div className="pt-navbar-heading">Nodes</div>
				</div>
				<div className="pt-navbar-group pt-align-right">
					<button
						className="pt-button pt-minimal pt-icon-refresh"
						onClick={this.sync}
					/>
				</div>
			</nav>
			{page}
		</div>;
	}
}
