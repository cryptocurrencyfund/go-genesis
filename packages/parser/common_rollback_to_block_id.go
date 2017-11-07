// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package parser

import (
	"database/sql"
	"fmt"

	"github.com/AplaProject/go-apla/packages/converter"
	"github.com/AplaProject/go-apla/packages/logging"
	"github.com/AplaProject/go-apla/packages/model"
)

// RollbackToBlockID rollbacks blocks till blockID
func (p *Parser) RollbackToBlockID(blockID int64) error {
	_, err := model.MarkVerifiedAndNotUsedTransactionsUnverified()
	if err != nil {
		logging.WriteSelectiveLog(err)
		return p.ErrInfo(err)
	}

	limit := 1000
	// roll back our blocks
	for {
		block := &model.Block{}
		blocks, err := block.GetBlocks(blockID, int32(limit))
		if err != nil {
			return p.ErrInfo(err)
		}
		if len(blocks) == 0 {
			break
		}
		fmt.Printf(`%s `, blocks[0].ID)
		for _, block := range blocks {
			// roll back our blocks to the block blockID
			err = BlockRollback(block.Data)
			if err != nil {
				return p.ErrInfo(err)
			}
		}
		blocks = blocks[:0]
	}
	block := &model.Block{}
	_, err = block.Get(blockID)
	if err != nil && err != sql.ErrNoRows {
		return p.ErrInfo(err)
	}
	data := block.Data
	converter.BytesShift(&data, 1)
	iblock := converter.BinToDecBytesShift(&data, 4)
	time := converter.BinToDecBytesShift(&data, 4)
	ecosystemID := converter.BinToDecBytesShift(&data, 4)
	keyID := converter.BinToDecBytesShift(&data, 8)
	nodePosition := converter.BinToDecBytesShift(&data, 1)
	ib := &model.InfoBlock{
		Hash:     block.Hash,
		BlockID:  iblock,
		Time:     time,
		EcosystemID: ecosystemID,
		KeyID: keyID,
		NodePosition:  nodePosition}
	err = ib.Update(p.DbTransaction)
	if err != nil {
		return p.ErrInfo(err)
	}
	err = model.UpdateConfig("my_block_id", converter.Int64ToStr(iblock))
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}
