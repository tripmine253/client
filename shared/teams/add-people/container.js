// @flow
import {connect} from 'react-redux'
import * as Creators from '../../actions/teams/creators'
import * as SearchCreators from '../../actions/search/creators'
import * as I from 'immutable'
import AddPeople from '.'
import {HeaderHoc} from '../../common-adapters'
import {compose, withPropsOnChange} from 'recompose'

import type {TypedState} from '../../constants/reducer'

const mapStateToProps = (state: TypedState, {routeProps}) => {
  const selectedUsersToAdd = state.entities.getIn(
    ['search', 'searchKeyToUserInputItemIds', 'addToTeamSearch'],
    I.List()
  )
  const usersAlreadyInTeam = state.entities.getIn([
    'teams',
    'teamNameToMemberUsernames',
    routeProps.get('teamname'),
  ])
  const tooManyUsers = selectedUsersToAdd.count() + usersAlreadyInTeam.count() >= 20
  return {
    name: routeProps.get('teamname'),
    tooManyUsers,
  }
}

const mapDispatchToProps = (dispatch: Dispatch, {navigateUp, routeProps}) => ({
  onClose: () => dispatch(navigateUp()),
  onAddPeople: (role: string) => {
    dispatch(Creators.addPeopleToTeam(routeProps.get('teamname'), role))
    dispatch(navigateUp())
    dispatch(Creators.getTeams())
    dispatch(SearchCreators.clearSearchResults('addToTeamSearch'))
    dispatch(SearchCreators.setUserInputItems('addToTeamSearch', []))
  },
})

export default compose(
  connect(mapStateToProps, mapDispatchToProps),
  compose(
    withPropsOnChange(['onExitSearch'], props => ({
      onCancel: () => props.onClose(),
      title: 'Add people',
    })),
    HeaderHoc
  )
)(AddPeople)
