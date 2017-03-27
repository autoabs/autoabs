/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	children?: JSX.Element;
}

export default class FixContainer extends React.Component<Props, void> {
	static childContextTypes = {};

	getChildContext() {
		return {};
	}

	render(): JSX.Element {
		return <div>
			{this.props.children}
		</div>;
	}
}
