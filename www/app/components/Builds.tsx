/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as BuildTypes from '../types/BuildTypes'
import BuildStore from '../stores/BuildStore';
import * as BuildUtils from '../utils/BuildUtils';
import Build from './Build';

interface Props {
	title: string;
}

interface State {
	builds: BuildTypes.Builds;
}

function getState(): State {
	return {
		builds: BuildStore.builds,
	};
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
	builds: {
		display: 'table',
		width: '100%',
	} as React.CSSProperties,
	buildsHeader: {
		display: 'table-row',
	} as React.CSSProperties,
	buildsHeaderField: {
		fontWeight: 800,
		display: 'table-cell',
	} as React.CSSProperties,
};

export default class Builds extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context);
		BuildUtils.init();
		this.state = getState();
	}

	componentDidMount(): void {
		BuildStore.addChangeListener(this._onChange);
	}

	componentWillUnmount(): void {
		BuildStore.removeChangeListener(this._onChange);
	}

	_onChange = (): void => {
		this.setState(getState());
	};

	render(): JSX.Element {
		let builds = this.state.builds;

		let buildsDom: JSX.Element[] = [];
		for (let key in builds) {
			buildsDom.push(<Build key={key} build={builds[key]}/>);
		}

		return <div>
			<paper-toolbar class="title">
				<div className="layout horizontal" style={css.headerBox}>
					<div className="flex" style={css.header}>
						{this.props.title}
					</div>
				</div>
			</paper-toolbar>
			<div style={css.builds}>
				<div style={css.buildsHeader}>
					<div style={css.buildsHeaderField}>ID</div>
					<div style={css.buildsHeaderField}>NAME</div>
					<div style={css.buildsHeaderField}>STATE</div>
					<div style={css.buildsHeaderField}>VERSION</div>
					<div style={css.buildsHeaderField}>RELEASE</div>
					<div style={css.buildsHeaderField}>REPO</div>
					<div style={css.buildsHeaderField}>ARCH</div>
					<div style={css.buildsHeaderField}>START</div>
					<div style={css.buildsHeaderField}>STOP</div>
				</div>
				{buildsDom}
			</div>
		</div>;
	}
}
