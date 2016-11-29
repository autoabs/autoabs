/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Card from 'material-ui/Card';
import CardText from 'material-ui/Card';
import CardActions from 'material-ui/Card';
import FlatButton from 'material-ui/FlatButton';
import IconButton from 'material-ui/IconButton';
import Dialog from 'material-ui/Dialog';
import Styles from '../Styles';
import * as BuildActions from '../actions/BuildActions';
import * as BuildTypes from '../types/BuildTypes';
import * as MiscUtils from '../utils/MiscUtils';

interface Props {
	build: BuildTypes.Build;
}

interface State {
	dialog: boolean;
}

const css = {
	card: {
		flex: '1 0 auto',
		minWidth: '300px',
		maxWidth: '450px',
		margin: '5px',
	} as React.CSSProperties,
	launch: {
		color: Styles.colors.color,
		margin: '1px 1px 0 0',
	} as React.CSSProperties,
	content: {
		padding: '10px 0 10px 10px',
	} as React.CSSProperties,
	name: {
		fontSize: '20px',
	} as React.CSSProperties,
	version: {
		fontSize: '14px',
		margin: '5px 0 0 7px',
		color: Styles.colors.fadeColor,
	} as React.CSSProperties,
	repo: {
		fontSize: '12px',
		marginTop: '7px',
		color: Styles.colors.fadeColor,
	} as React.CSSProperties,
	actions: {
		display: 'flex',
		flexDirection: 'row',
		padding: '5px 0',
		justifyContent: 'center',
	} as React.CSSProperties,
	skip: {
		color: Styles.colors.yellow400,
	} as React.CSSProperties,
	resume: {
		color: Styles.colors.green500,
	} as React.CSSProperties,
	retry: {
		color: Styles.colors.blue500,
	} as React.CSSProperties,
	remove: {
		color: Styles.colors.red500,
	} as React.CSSProperties,
};

export default class Build extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context);
		this.state = {
			dialog: false,
		};
	}

	openDialog = (): void => {
		this.setState({
			dialog: true,
		});
	}

	closeDialog = (): void => {
		this.setState({
			dialog: false,
		});
	}

	onRemove = (): void => {
		BuildActions.remove(this.props.build.id);
	}

	render(): JSX.Element {
		let build = this.props.build;

		let start = '-';
		if (build.start !== '0001-01-01T00:00:00Z') {
			start = MiscUtils.formatDate(new Date(build.start));
		}
		let stop = '-';
		if (build.stop !== '0001-01-01T00:00:00Z') {
			stop = MiscUtils.formatDate(new Date(build.stop));
		}

		let actions: JSX.Element[];

		switch (build.state) {
			case "building":
				actions = [
					<FlatButton style={css.skip} label="Skip"/>,
				];
				break;
			case "pending":
				actions = [
					<FlatButton style={css.skip} label="Skip"/>,
					<FlatButton style={css.remove}
						label="Remove" onTouchTap={this.onRemove}/>,
				];
				break;
			case "failed":
				actions = [
					<FlatButton style={css.retry} label="Retry"/>,
					<FlatButton style={css.remove}
						label="Remove" onTouchTap={this.onRemove}/>,
				];
				break;
			case "completed":
				actions = [
					<FlatButton style={css.retry} label="Retry"/>,
					<FlatButton style={css.remove}
						label="Remove" onTouchTap={this.onRemove}/>,
				];
				break;
			case "skipped":
				actions = [
					<FlatButton style={css.retry} label="Retry"/>,
					<FlatButton style={css.remove}
						label="Remove" onTouchTap={this.onRemove}/>,
				];
				break;
		}

		let dialogActions = [
			<FlatButton
				label="Close"
				primary={true}
				onTouchTap={this.closeDialog}
			/>,
		];

		return <Card style={css.card}>
			<CardText>
				<div className="layout horizontal">
					<div style={css.content} className="card-content flex">
						<div className="layout vertical">
							<div className="layout horizontal">
								<div style={css.name}>{build.name}</div>
								<div style={css.version}>{build.version}-{build.release}</div>
							</div>
							<div style={css.repo}>{build.repo} - {build.arch}</div>
						</div>
					</div>
					<div>
						<IconButton style={css.launch}
								onClick={this.openDialog}>
							<i className="material-icons">flip_to_front</i>
						</IconButton>
					</div>
				</div>
			</CardText>
			<CardActions style={css.actions}>
				{actions}
			</CardActions>
			<Dialog
				title={`Builds Logs - ${build.name}`}
				modal={true}
				actions={dialogActions}
				open={this.state.dialog}
			>Test dialog</Dialog>
		</Card>;
	}
}
