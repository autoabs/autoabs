/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import LoadingStore from '../stores/LoadingStore';

interface Props {
	style?: React.CSSProperties;
	size?: string;
	intent?: Blueprint.Intent,
}

interface State {
	loading: boolean;
}

export default class Loading extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context);
		this.state = {
			loading: LoadingStore.loading,
		};
	}

	componentDidMount(): void {
		LoadingStore.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		LoadingStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			loading: LoadingStore.loading,
		});
	}

	render(): JSX.Element {
		if (!this.state.loading) {
			return null;
		}

		let className = '';
		if (this.props.size) {
			className = 'pt-' + this.props.size;
		}

		return <div style={this.props.style}>
			<Blueprint.Spinner
				className={className}
				intent={this.props.intent}
			/>
		</div>;
	}
}
