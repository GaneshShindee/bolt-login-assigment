import type { SavedDetails } from '../api'
import { formatAddress, labelIcon } from '../format'
import type { Choice } from '../hooks/useDeliveryChoice'

type Props = {
  saved: SavedDetails[]
  selected: Choice
  editing: boolean
  compact: boolean
  newAddress: { typed: boolean; preview?: string }
  onSelect: (choice: Choice) => void
  onEdit: () => void
}

export default function SavedDetailsPicker({ saved, selected, editing, compact, newAddress, onSelect, onEdit }: Props) {
  return (
    <div className={`saved-grid${compact ? ' compact' : ''}`} role="radiogroup" aria-label="Delivery address">
      {saved.map((d, i) => {
        const isSelected = selected === i
        return (
          <label key={i} className={`saved-card${isSelected ? ' selected' : ''}`}>
            <input
              type="radio"
              name="saved-details"
              className="visually-hidden"
              checked={isSelected}
              onChange={() => onSelect(i)}
            />
            <span className="saved-card-head">
              <span className="saved-label">
                {labelIcon(d.address.label)} {d.address.label}
              </span>
              {isSelected &&
                (editing ? (
                  <span className="muted small">Editing</span>
                ) : (
                  <button type="button" className="btn link small inline" onClick={onEdit}>
                    Edit
                  </button>
                ))}
            </span>
            <span className="saved-address">{formatAddress(d.address)}</span>
            <span className="muted small">📞 {d.phone}</span>
          </label>
        )
      })}
      <label className={`saved-card new${selected === 'new' ? ' selected' : ''}`}>
        <input
          type="radio"
          name="saved-details"
          className="visually-hidden"
          checked={selected === 'new'}
          onChange={() => onSelect('new')}
        />
        <span className="saved-new">{newAddress.typed ? '✎ Your new address' : '+ New address'}</span>
        {newAddress.preview && <span className="saved-address">{newAddress.preview}</span>}
      </label>
    </div>
  )
}
