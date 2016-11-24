/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Builds from './Builds';

interface Props {
	title: string;
}

const css = {
	headerBox: {
		color: '#fff',
		width: '100%',
	} as React.CSSProperties,
	header: {
		margin: '4px',
		fontSize: '24px',
	} as React.CSSProperties,
};

export default class Main extends React.Component<Props, null> {
	render(): JSX.Element {
		return <div>
			<paper-toolbar class="title">
				<div className="layout horizontal" style={css.headerBox}>
					<div className="flex" style={css.header}>
						{this.props.title}
					</div>
				</div>
			</paper-toolbar>
			<Builds/>
		</div>;
	}
}
