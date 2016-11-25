/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Builds from './Builds';

interface Props {
	title: string;
}

const css = {
	header: {
		color: '#fff',
		marginLeft: '4px',
		fontSize: '24px',
	} as React.CSSProperties,
};

export default class Main extends React.Component<Props, null> {
	render(): JSX.Element {
		return <div>
			<paper-toolbar>
				<div className="title" style={css.header}>{this.props.title}</div>
			</paper-toolbar>
			<Builds/>
		</div>;
	}
}
