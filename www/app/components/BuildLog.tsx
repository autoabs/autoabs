/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import Loading from './Loading';
import * as BuildLogActions from '../actions/BuildLogActions';
import BuildStore from '../stores/BuildStore';
import BuildLogStore from '../stores/BuildLogStore';

interface State {
	id: string;
	name: string;
	output: string;
}

function getState(): State {
	let id = BuildLogStore.id;
	let name = '';

	if (id) {
		let build = BuildStore.build(id);
		if (build) {
			name = build.name;
		}
	}

	return {
		id: id,
		name: name,
		output: BuildLogStore.output,
	};
}

const css = {
	loading: {
		float: "left",
	} as React.CSSProperties,
	buildLog: {
		top: '20px',
		width: 'calc(100% - 40px)',
		height: 'calc(100% - 40px)',
	} as React.CSSProperties,
	buildLogOutput: {
		fontSize: '10px',
		width: '100%',
		height: 'calc(100% - 130px)',
		overflow: 'scroll',
		padding: '2px 6px',
		backgroundColor: 'rgba(0, 0, 0, 0.2)',
	} as React.CSSProperties,
};

export default class Build extends React.Component<void, State> {
	constructor(props: void, context: any) {
		super(props, context);
		this.state = getState();
	}

	componentDidMount(): void {
		BuildStore.addChangeListener(this.onChange);
		BuildLogStore.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		BuildStore.removeChangeListener(this.onChange);
		BuildLogStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState(getState());
	}

	closeDialog = (): void => {
		BuildLogActions.close();
	}

	render(): JSX.Element {
		return <Blueprint.Dialog
			title={`Builds Logs - ${this.state.name}`}
			style={css.buildLog}
			isOpen={!!this.state.id}
			onClose={this.closeDialog}
			canOutsideClickClose={false}
		>
			<div className="pt-dialog-body">
					<pre style={css.buildLogOutput}>
						{this.state.output}
					</pre>
			</div>
			<div className="pt-dialog-footer">
				<Loading size="small" style={css.loading}/>
				<div className="pt-dialog-footer-actions">
					<button type="button"
						className="pt-button"
						onClick={this.closeDialog}
					>Close</button>
				</div>
			</div>
		</Blueprint.Dialog>;
	}
}
